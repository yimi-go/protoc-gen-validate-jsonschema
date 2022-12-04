.PHONY: release.tag
release.tag: tools.verify.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
	@build/scripts/ensure_tag.sh