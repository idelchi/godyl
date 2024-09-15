1. Observe the established pattern in internal/cli/{auth,cache,config}/\*\*, in particular the subcommands (for example internal/cli//auth/dump/command.go, internal/cli//auth/dump/logic.go and so on)
   - Package level and function level documentation and wording is consistent.
   - Usage of common.Input to pass to run(), where run() does all the work.
   - Short, Long, Example fields are used (where they make sense)

Apply the same pattern to (where it is applicable):

    - internal/cli/{download,dump,install,status,update,validate,version}

ONLY focus on what was described in (1) and do not change anything else.
Important is CONSISTENCY with the existing pattern. Do not STRICTLY adhere to the consistency - for example, many subcommands need `embedded` to be passed,
as such, do not blatantly remove it, but do make sure it is added to the common.Input struct and used in run().

Make sure each command.go and logic.go have the same wording to the package level and function level comments, that is, if there's an existing one that does not adhere to the pattern, change it to match the others.
