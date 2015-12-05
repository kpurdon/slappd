# Slack Untappd (WORK IN PROGRESS)

A [Slack](https://slack.com/) [slash command server](https://api.slack.com/slash-commands) for searching for beer information via the [Untappd Beer Search API](https://untappd.com/api/docs#userbeers).

## Basic Usage

1. Register for [Untappd API Credentials](https://untappd.com/api/register?register=new)
2. Get the server running somewhere on a public host (TODO)
3. Configure a Slack Slash Command. [See Slack Docs](https://api.slack.com/slash-commands)

The following environment variables must be set:

* `SLACK_TOKEN` - The token assigned to the slash command integration
* `UNTAPPD_CLIENT_ID` - Your Untappd API client ID
* `UNTAPPD_CLIENT_SECRET` - Your Untappd API client secret token

### Request

The server expects a POST request from Slack with the following form values:

* token
* text

These values should be present in the default slack POST in addition to other non-used values.

### Response

The server will do a lookup on the Untappd API for the given `text` in the Slack POST. If no result is found an `ephemeral` (only displayed to requesting user) message will be returned stating that no results were found. If results are found an `in_channel` (displayed to anyone in the channel) message will be returned.
