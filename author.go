package main
import "strings"

func getDistinguished(distinguished string, isSubmitter bool) (ua userAttrs) {

	ua = userAttrs{
		Submitter: isSubmitter,
		Moderator: strings.Contains(distinguished, "moderator"),
		Admin:     strings.Contains(distinguished, "admin"),
	}
	ln := numTrues(ua.Moderator, ua.Admin, ua.Submitter)
	if ln == 0 {
		ua.Letters = ""
		return
	}
	tmp := make([]string, ln)
	i := 0
	if ua.Moderator {
		tmp[i] = "<a href=\"#\" class=\"moderator\" title=\"moderator of this subreddit, speaking officially\">M</a>"
		i++
	}
	if ua.Admin {
		tmp[i] = "<a href=\"#\" class=\"admin\" title=\"Reddit Administrator\">A</a>"
		i++
	}
	if ua.Submitter {
		tmp[i] = "<a href=\"#\" class=\"submitter\" title=\"submitter\">S</a>"
		i++
	}
	ua.Letters = "[" + strings.Join(tmp, ",") + "]"
	return
}

func isMine(author string) bool {
	return strings.
		EqualFold(client.Username, author)
}