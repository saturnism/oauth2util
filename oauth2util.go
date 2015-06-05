/*
 * Copyright 2015 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"bufio"
	"fmt"
	"os"

	"encoding/json"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/internal"

	"github.com/codegangsta/cli"
)

func auth(conf *oauth2.Config) (*oauth2.Token, error) {
	url := conf.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter the authorization code: ")
	code, _ := reader.ReadString('\n')
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "oauth2util"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "client_id",
			Usage:  "Client ID",
			EnvVar: "OAUTH2_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "client_secret",
			Usage:  "Client Secret",
			EnvVar: "OAUTH2_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "auth_url",
			Value:  "https://accounts.google.com/o/oauth2/auth",
			Usage:  "OAuth 2.0 Authorization URL",
			EnvVar: "OAUTH2_AUTH_URL",
		},
		cli.StringFlag{
			Name:   "token_url",
			Value:  "https://accounts.google.com/o/oauth2/token",
			Usage:  "OAuth 2.0 Token URL",
			EnvVar: "OAUTH2_TOKEN_URL",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "exchange",
			Aliases: []string{"e"},
			Usage:   "Exchange refresh token for access token",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "refresh_token, t",
					Usage:  "Refresh token",
					EnvVar: "OAUTH2_REFRESH_TOKEN",
				},
			},
			Action: func(c *cli.Context) {
				conf := &ExtendedOAuth2Config{
					ClientID:     c.GlobalString("client_id"),
					ClientSecret: c.GlobalString("client_secret"),
					Endpoint: oauth2.Endpoint{
						AuthURL:  c.GlobalString("auth_url"),
						TokenURL: c.GlobalString("token_url"),
					},
				}
				token, err := conf.ExchangeWithRefreshToken(oauth2.NoContext, c.String("refresh_token"))
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				prettyPrint(&token)
			},
		},
		{
			Name:    "auth",
			Aliases: []string{"a"},
			Usage:   "Authorize using OAuth 2.0",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "redirect_url, r",
					Value:  "urn:ietf:wg:oauth:2.0:oob",
					Usage:  "Redirect URL",
					EnvVar: "OAUTH2_REDIRECT_URL",
				},
				cli.StringSliceFlag{
					Name:   "scopes, s",
					Usage:  "OAuth 2.0 Scopes",
					EnvVar: "OAUTH2_SCOPES",
				},
			},
			Action: func(c *cli.Context) {
				conf := &oauth2.Config{
					ClientID:     c.GlobalString("client_id"),
					ClientSecret: c.GlobalString("client_secret"),
					RedirectURL:  c.String("redirect_url"),
					Scopes:       c.StringSlice("scopes"),
					Endpoint: oauth2.Endpoint{
						AuthURL:  c.GlobalString("auth_url"),
						TokenURL: c.GlobalString("token_url"),
					},
				}
				token, err := auth(conf)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				prettyPrint(&token)
			},
		},
	}
	app.Usage = "Easy command line access to generate OAuth 2.0 credentials"

	app.Run(os.Args)
}

type ExtendedOAuth2Config oauth2.Config

func (c *ExtendedOAuth2Config) ExchangeWithRefreshToken(ctx context.Context, refresh_token string) (*oauth2.Token, error) {
	v := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh_token},
	}

	t, err := internal.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.Endpoint.TokenURL, v)
	if err != nil {
		return nil, err
	}
	return tokenFromInternal(t), nil
}

func tokenFromInternal(t *internal.Token) *oauth2.Token {
	if t == nil {
		return nil
	}
	return &oauth2.Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		// raw is missing, but doesn't seem necessary here...
		// raw:        t.Raw,
	}
}

func prettyPrint(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
