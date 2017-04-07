// tagmail - tag your folders as your tags as your folders
//
// This script makes 2 huge assumptions and is why we use Go: that you have
// parallel I/O bandwidth (i.e. flash) and that you have parallel CPU bandwidth
// (i.e. multiple cores). We'll do our best to be super performant over large
// mailboxes, and maybe ever rewrite this in C later.. but for now we assume the
// most important thing is to issue the reads to disk as fast as possible, then
// slowly optimize from there.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	dry = flag.Bool("dry", false,
		"Run without actually executing modifications")

	maildir = flag.String("mailroot", os.Getenv("MAIL"),
		"Root of tree of maildirs to syncronize tags")

	inbox = flag.String("inbox", "INBOX", "INBOX directory")
	trash = flag.String("trash", "Trash", "Trash directory")
	sent  = flag.String("sent", "Sent Items", "Sent directory")

	tag_unread = flag.String("tag-unread", "unread", "unread tag name")
	tag_new    = flag.String("tag-new", "new", "new tag name")

	notmuch_root   = flag.String("notmuch-root", notmuchRoot(), "Notmuch Root Location")
	nomultiaccount = flag.Bool("nomultiaccount", false, "Run without support for multiple accounts")
)

func main() {
	flag.Parse()

	// MAILBOXES_FULL_PATHS="$(echo "$(find $MAILDIR_ACCOUNT_ROOT -name "cur" -type d -exec dirname '{}' \;)" | sort;)"
	//  # | sed "s/^$MAILDIR\///" | sort
	mailboxes, err := findMailboxes(*maildir)
	if err != nil {
		panic(err)
	}

	for _, v := range mailboxes {
		fmt.Printf("mailbox: %s\n", v)
	}

}

func findMailboxes(mailroot string) ([]string, error) {
	mailboxnames := make([]string, 0)
	err := filepath.Walk(mailroot, func(path string, f os.FileInfo, err error) error {
		path = strings.Trim(strings.TrimPrefix(path, mailroot), "/")
		if !strings.HasPrefix(path, ".") && path != "" {
			//			return fmt.Errorf("dotted path")
			mailboxnames = append(mailboxnames, path)
		}

		return nil
	})
	return mailboxnames, err
}

// This takes no pathoptions yet, but can in the future
// Notmuch_Tag_From_Full_Path ()
func pathToTag(mailroot, path string) string {
	p := strings.TrimPrefix(path, mailroot)
	p = strings.Trim(p, "/")
	if !nomultiaccount {
		slash_idx = strings.Index(p, "/")
		p = p[:slash_idx]
	}
	return p
}

// Notmuch_Folder_From_Full_Path
func pathToNotmuchFolder(mailroot, path string) string {
	p := strings.TrimPrefix(path, mailroot)
	p = strings.Trim(p, "/")
}

func notmuchRoot() string {
	out, err := exec.Command("notmuch", "config", "list").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get notmuch config\n")
		panic(err)
	}

	db_index := bytes.Index(out, []byte("database.path="))
	if db_index == -1 {
		panic("Could not find database path?!")
	}

	newline := bytes.Index(out[db_index:], []byte("\n"))
	return string(out[db_index+len("database.path=") : newline])
}

// Maildir_Account_Folder_From_Full_Path ()
func pathToAcct(mailroot, path string) string {
	p := strings.TrimPrefix(path, mailrooc)
	p = strings.Trim(p, "/")
	slash_idx = strings.Index(p, "/")
	if slash_idx == -1 {
		return p
	}

	// make sure this doesn't return the slash?!
	return p[:slash_idx]
}

func st_tagAdded_maildirMissing(mailroot, mailpath string) {
	notmuchFolder :=  pathToNotmuchFolder(mailroot, mailpath)
	notmuchTag := pathToTag(mailroot, mailpath)
	notmuchAcct := pathToAcct(mailroot, mailpath)
	cmd := exec.Command("notmuch",
		"search", "--output=messages",
		"tag:" + notmuchTag,
		"path:" + notmuchAcct + "/**",
		"NOT", "folder:" + notmuchFolder,
		"NOT", "folder:" + *trash,
		"NOT", "tag:new")

	msgs, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	notmuch search --output
}
func st_tagRemoved_maildirExtra() {}

func st_maildirAdded_tagMissing() {}
func st_maildirRemoved_tagExtra() {}

// # executed prior to 'notmuch new'
// Notmuch_State_To_Maildir__Move_To_Maildir ()
// {
// # Scenario:
// #
// # NOTMUCH STATE (per message):
// # Number of Notmuch Tags > Number of Notmuch Folders
// #
// # MAILDIR STATE:
// # No change from previous state.
// #
// # Tags have been added to a message in a virtual folder (in the notmuch db).
// # The number of folders associated with a message has not been changed in
// # the notmuch db. This indicates that we need to copy the message to a
// # new maildir. After the next 'notmuch new' db update, the tags/folders
// # should thus be at parity again.

//     # We are running this function prior to the remove function below
//     # but there is still the edge case wherein the user has manually
//     # deleted a mail message in mutt (better to do all this with tags
//     # and virtual folders, but let's accommodate).

//     for THIS_MESSAGE_ID in $THESE_MESSAGE_IDS_TO_COPY; do

//         local THIS_MESSAGE_ALL_SOURCE_PATHS="$(notmuch search --output=files "$THIS_MESSAGE_ID")"

//         local FOUND=false

//         while read line; do

//             local THIS_MESSAGE_SOURCE_PATH="$line"

//             if [[ -e "$THIS_MESSAGE_SOURCE_PATH" ]]; then
//                 FOUND=true
//                 break
//             fi

//         done <<< "$THIS_MESSAGE_ALL_SOURCE_PATHS"

//         if $FOUND; then

//             if $RUNCMD "cp \"$THIS_MESSAGE_SOURCE_PATH\" \"$THIS_MAILDIR_FULL_PATH/cur\""; then
//                 echo -n "Copied message with new tag to"
//                 echo " $(Maildir_Account_Folder_From_Full_Path "$THIS_MAILDIR_FULL_PATH")"
//             else
//                 echo -e "\nWARNING: Failed to copy mail file (unknown error):"
//                 echo -e "SOURCE:  \"$THIS_MESSAGE_SOURCE_PATH\"\nDESTINATION\"$THIS_MAILDIR_FULL_PATH/cur\"\n"
//             fi

//         else

//             echo -e "\nWARNING: Failed to copy mail file (no valid source paths!):"
//             echo "ID: $THIS_MESSAGE_ID"
//             echo "NOTMUCH FOLDER: $THIS_NOTMUCH_FOLDER"
//             echo -e "DESTINATION MAILDIR: $THIS_MAILDIR_FULL_PATH\n"
//         fi

//     done

// }

// Notmuch_State_To_Maildir__Remove_From_Maildir ()
// {
// # Scenario:
// #
// # NOTMUCH STATE (per message):
// # Number of Notmuch Tags < Number of Notmuch Folders
// #
// # MAILDIR STATE:
// # No change from previous state.
// #
// # Tags have been removed from a message in a virtual folder (and thus
// # in the notmuch db). The number of folders associated with a message
// # has of course not yet changed. We need to remove the messages from
// # maildir folders from which it has been untagged.

//     local THIS_MAILDIR_FULL_PATH="$1"
//     local THIS_NOTMUCH_FOLDER="$(Notmuch_Folder_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THIS_NOTMUCH_TAG="$(Notmuch_Tag_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THESE_MESSAGE_IDS_TO_REMOVE="$(\
//         notmuch search --output=messages\
//         folder:"$THIS_NOTMUCH_FOLDER" \
//         NOT tag:"$THIS_NOTMUCH_TAG" \
//         NOT tag:"$NEW_TAG")"

//     for THIS_MESSAGE_ID in $THESE_MESSAGE_IDS_TO_REMOVE; do

//         local THIS_MESSAGE_PATH="$(notmuch search --output=files "$THIS_MESSAGE_ID" | \
//             grep -e "^$THIS_MAILDIR_FULL_PATH")"

//         if [[ -e "$THIS_MESSAGE_PATH" ]]; then

//             if $RUNCMD "rm \"$THIS_MESSAGE_PATH\""; then
//                 echo -n "Removed untagged message from"
//                 echo " $(Maildir_Account_Folder_From_Full_Path "$THIS_MAILDIR_FULL_PATH")"
//             else
//                 echo -e "\nWARNING: Failed to remove mail file (unknown error):"
//                 echo "ID:$THIS_MESSAGE_ID"
//                 echo "FOLDER:$THIS_NOTMUCH_FOLDER"
//                 echo -e "MESSAGE PATH:$THIS_MESSAGE_PATH\n"
//             fi

//         else

//             echo -e "\nWARNING: Unable to remove missing mail file:"
//             echo "ID:$THIS_MESSAGE_ID"
//             echo "FOLDER:$THIS_NOTMUCH_FOLDER"
//             echo -e "MESSAGE PATH:$THIS_MESSAGE_PATH\n"
//         fi

//     done

// }

// # ----------------------------------------------------------------------
// # SYNC Notmuch DB Sync Functions
// # ----------------------------------------------------------------------

// Notmuch_Update ()
// {
//     $RUNCMD "notmuch new";
//     # FIXME pass TAG_SCRIPT as an argument
//     if [[ -e "$TAG_SCRIPT" ]]; then
//         $RUNCMD $TAG_SCRIPT
//     fi
// }

// # ----------------------------------------------------------------------
// # POST Notmuch DB Sync Functions
// # ----------------------------------------------------------------------
// # executed after 'notmuch new' (otherwise the notmuch state looks the
// # same as the states above)

// Maildir_State_To_Notmuch__Add_Tags_To_Notmuch ()
// {
// # Scenario:
// #
// # NOTMUCH STATE (per message):
// # Number of Notmuch Tags < Number of Notmuch Folders
// #
// # MAILDIR STATE:
// # Message in a new folder (either via CLI/mutt copy, move or incoming sync)
// #
// # A message is in a "physical" maildir directory but does not have a
// # corresponding notmuch tag. For example:
// #
// #     ~/mail/INBOX/message123 should have a tag "INBOX"
// #
// # We process all mails in each maildir directory (mailbox) and add tags
// # as required.

//     local THIS_MAILDIR_FULL_PATH="$1"
//     local THIS_NOTMUCH_FOLDER="$(Notmuch_Folder_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THIS_NOTMUCH_TAG="$(Notmuch_Tag_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THIS_NOTMUCH_QUERY="folder:\"$THIS_NOTMUCH_FOLDER\" NOT tag:\"$THIS_NOTMUCH_TAG\""
//     local THIS_COUNT="$(notmuch count $THIS_NOTMUCH_QUERY)"

//     $DRYRUN || notmuch tag +"$THIS_NOTMUCH_TAG" -- $THIS_NOTMUCH_QUERY
//     [[ $THIS_COUNT > 0 ]] && echo "Tagged $THIS_COUNT messages with \"$THIS_NOTMUCH_TAG\"" || true

// }

// Maildir_State_To_Notmuch__Remove_Tags_From_Notmuch ()
// {
// # Scenario:
// #
// # NOTMUCH STATE (per message):
// # Number of Notmuch Tags > Number of Notmuch Folders
// #
// # MAILDIR STATE:
// # Message removed from folder, either via rm, mutt delete, or offlineimap sync
// #
// # A message has been removed from a maildir directory. Notmuch is aware of
// # this (this should only be checked/run after a 'notmuch new' update).
// # However, we still have the "old" tag on the message.
// #
// # We skip the trash since we might want to restore those in future?

//     local THIS_MAILDIR_FULL_PATH="$1"
//     local THIS_NOTMUCH_FOLDER="$(Notmuch_Folder_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THIS_NOTMUCH_TAG="$(Notmuch_Tag_From_Full_Path $THIS_MAILDIR_FULL_PATH)"
//     local THIS_ACCOUNT_FOLDER="$(dirname $THIS_NOTMUCH_FOLDER)"
//     local THIS_NOTMUCH_QUERY="tag:\"$THIS_NOTMUCH_TAG\"\
//                               path:\"$THIS_ACCOUNT_FOLDER/**\"\
//                               NOT folder:\"$THIS_NOTMUCH_FOLDER\" \
//                               NOT folder:\"$TRASH\""
//     local THIS_COUNT="$(notmuch count $THIS_NOTMUCH_QUERY)"

//     $DRYRUN || notmuch tag -"$THIS_NOTMUCH_TAG" -- $THIS_NOTMUCH_QUERY
//     [[ $THIS_COUNT > 0 ]] && echo "Untagged $THIS_COUNT messages, removed \"$THIS_NOTMUCH_TAG\"" || true

// }

// # ----------------------------------------------------------------------
// # CLEANUP Functions
// # ----------------------------------------------------------------------

// Notmuch_Cleanup ()
// {

//     # anything in sent mail should have the unread flag removed
//     $RUNCMD "notmuch tag -\"$UNREAD_TAG\" -- folder:\"$SENT\""

//     # FIXME pass CLEANUP_SCRIPT as an argument
//     if [[ -e "$CLEANUP_SCRIPT" ]]; then
//         $RUNCMD $CLEANUP_SCRIPT
//     fi

//     # remove "$NEW_TAG" tags, optionally converting to "$UNREAD_TAG"
//     case $MAKE_NEW_UNREAD in
//         true|TRUE|yes|YES|y|Y)
//             $RUNCMD "notmuch tag -\"$NEW_TAG\" +\"$UNREAD_TAG\" -- tag:\"$NEW_TAG\"" ;;
//         *)  $RUNCMD "notmuch tag -\"$NEW_TAG\" -- tag:\"$NEW_TAG\"" ;;
//     esac

// }

// # ----------------------------------------------------------------------
// # ----------------------------------------------------------------------
// # MAIN
// # ----------------------------------------------------------------------
// # ----------------------------------------------------------------------

// echo -e "\n----------------------------------------------------------------------"
// echo "$(basename $0) ${SUBCMD}-sync hook ${DRYRUN_MSG:-}"
// echo "----------------------------------------------------------------------"
// echo "NOTMUCH ROOT: $NOTMUCH_ROOT"
// echo "ACCOUNT ROOT: $MAILDIR_ACCOUNT_ROOT"

// # Review the notmuch database state and sync up any changes first
// # (e.g. any retagged messages that need refiling)
// if [ "$SUBCMD" == "pre" ]; then
//     for MAILBOX_FULL_PATH in $MAILBOXES_FULL_PATHS; do
//         Notmuch_State_To_Maildir__Move_To_Maildir $MAILBOX_FULL_PATH
//         Notmuch_State_To_Maildir__Remove_From_Maildir $MAILBOX_FULL_PATH
//     done
// fi

// # Update the notmuch database to reflect the changes we just made,
// # if any (so it can find the new messages)
// if [ "$SUBCMD" == "post" ]; then
//     Notmuch_Update

//     for MAILBOX_FULL_PATH in $MAILBOXES_FULL_PATHS; do
//         Maildir_State_To_Notmuch__Add_Tags_To_Notmuch $MAILBOX_FULL_PATH
//         Maildir_State_To_Notmuch__Remove_Tags_From_Notmuch $MAILBOX_FULL_PATH
//     done

//     Notmuch_Cleanup
// fi

// echo -e "maildir-notmuch-sync complete ----------------------------------------\n"
