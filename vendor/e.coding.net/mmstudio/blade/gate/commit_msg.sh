#!/bin/sh
brew install node
echo "node_modules\npackage-lock.json\n">>.gitignore
echo '{
  "version": "1.0.0",
  "config": {
    "commitizen": {
      "path": "./node_modules/cz-conventional-changelog"
    },
    "ghooks": {
      "commit-msg": "validate-commit-msg"
    }
  },
  "scripts": {
    "all-log": "conventional-changelog -p angular -i ALLCHANGELOG.md -s -r 0",
    "changelog": "conventional-changelog -p angular -i CHANGELOG.md -s",
	"release": "standard-version",
	"postrelease": "git push --follow-tags origin $(git rev-parse --symbolic-full-name --abbrev-ref HEAD)"
  }
}' |jq . >package.json

npm install -g commitizen
npm install ghooks validate-commit-msg --save-dev
npm install -g conventional-changelog-cli
commitizen init cz-conventional-changelog -save -save-exact --force
npm install standard-version --save-dev
npm install -g i5ting_toc
sed -i '' 's/is_auto_number:true/is_auto_number:false/g' node_modules/i5ting_toc/vendor/toc_conf.js
echo '{
  "types": [
    {"type":"feat","section":"Features"},
    {"type":"fix","section":"Bug Fixes"},
    {"type": "refactor", "section": "Refactor"},
    {"type": "docs", "section": "Docs","hidden": true},
    {"type": "style", "section": "Style", "hidden": true},
    {"type":"test","section":"Tests", "hidden": true},
    {"type":"build","section":"Build System", "hidden": true},
    {"type":"ci","hidden":true},
    {"type":"chore","hidden":true}
  ]
}' >.versionrc

echo 'install succ';
