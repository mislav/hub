" Vim syntax file
" Language: Hub Pull Request
" Maintainer: Derek Sifford <dereksifford@gmail.com>
" Filenames: *.git/PULLREQ_EDITMSG
" Latest Revision: 2018 Oct 30

if exists('b:current_syntax')
    finish
endif

syn case match

syn include @Markdown syntax/markdown.vim

syn match pullreqBlank         contained                  "^.*"        contains=@Spell
syn match pullreqOverflow      contained                  ".*"         contains=@Spell
syn match pullreqSummary       contained                  "^.\{0,50\}" contains=@Spell nextgroup=pullreqOverflow
syn match pullreqMetaHeader    contained                  "^Changes:"
syn match pullreqSha           contained                  "^[a-z0-9]\{7\}\ze ("        nextgroup=pullreqCommitMeta
syn match pullreqCommitMeta    contained skipnl skipwhite ".*"                         nextgroup=pullreqCommitMessage
syn match pullreqCommitMessage contained                  "^\s*\zs.*"
syn match pullreqBranchInfo    contained                  "\S\+:\S\+"

syn region pullreqBranchInfoLine contained transparent start="^Requesting a pull" end="$"             contains=pullreqBranchInfo
syn region pullreqMessage        keepend               start="^."                 end="^\ze# [-]* >8" contains=@Markdown,@Spell                                   nextgroup=pullreqMetadata
syn region pullreqMetadata       fold                  start="^# [-]* >8 [-]*$"   end="\%$"           contains=pullreqMetaHeader,pullreqSha,pullreqBranchInfoLine
syn match  pullreqFirstLine      skipnl                "\%^[^#].*"                                    contains=pullreqSummary                                     nextgroup=pullreqBlank

hi def link pullreqBlank         Error
hi def link pullreqBranchInfo    Keyword
hi def link pullreqCommitMessage String
hi def link pullreqMetaHeader    htmlH1
hi def link pullreqMetadata      Comment
hi def link pullreqSha           Constant
hi def link pullreqSummary       Keyword

let b:current_syntax = 'pullrequest'
