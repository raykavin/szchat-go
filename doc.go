// Package szchat is a Go SDK for the SZChat Chat Center API
// (see the official SZChat API documentation, available under the
// "/docs/pt-br" path of your tenant's SZChat host).
//
// A Client authenticates lazily: the first request triggers a login using
// the credentials passed to NewClient, and the bearer token is refreshed
// automatically whenever a request comes back unauthorized.
//
// baseURL is tenant-specific (e.g. "https://your-tenant.sz.chat") and must
// not include the API version path: NewClient appends "/api/v4" to it
// automatically to keep internal compatibility with the API version this
// library targets.
//
//	client, err := szchat.NewClient(baseURL, email, password)
//	if err != nil {
//		return err
//	}
//
//	contacts, err := client.ContactAPI.List(ctx, szchat.ContactListFilter{})
//	if err != nil {
//		if szchat.IsRateLimited(err) {
//			// back off and retry later
//		}
//		return err
//	}
//
//	for _, contact := range contacts.Data {
//		fmt.Println(contact.Name)
//	}
//
// Every resource lives under its own Client field (ContactAPI, AgentAPI,
// TeamAPI, ChannelAPI, and so on), each wrapping a single business domain of
// the SZChat API.
package szchat
