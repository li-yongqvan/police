# AI Forum Context

AI Forum is a student-facing AI community whose mobile forum experience should support quick reading, topic scanning, and restrained content-first interaction.

## Language

**Mobile Quick Scan**:
The mobile forum home experience prioritizes quickly scanning topics over preserving desktop-style welcome, search, and filter density.
_Avoid_: Desktop parity, portal-style mobile home

**Compact Action Entry**:
When mobile space is tight, an action entry may use a clear icon first. If the icon alone is not recognizable enough, it may use one representative Chinese character while preserving the action's semantic label.
_Avoid_: Unlabeled color block, cramped full-label button

**Header Search**:
On the mobile forum home, search remains a header action rather than a large in-page search box. The header search entry must be easy to tap and must open a usable search surface.
_Avoid_: In-page mobile search hero

**Mobile Topic Heat Anchor**:
On the mobile forum topic list, the reply count remains a right-side scan anchor for discussion heat. It should be easier to tap than default compact text, while remaining visually secondary to the topic title.
_Avoid_: Tiny reply link, metadata-only reply count

**Quiet Topic Metadata**:
On the mobile forum topic list, category and activity time are supporting metadata. They should stay quiet and orderly rather than becoming prominent pill controls.
_Avoid_: Prominent metadata pills, category-as-primary-action

**Home-First Mobile Theme Fix**:
The first MUI-02 implementation batch focuses on the mobile forum home topic list rather than topic detail or composer surfaces.
_Avoid_: Detail-page scope creep

**Mobile UI Hard Acceptance**:
Mobile UI work is accepted against concrete viewport, overflow, touch target, text wrapping, and desktop regression checks rather than broad visual preference.
_Avoid_: Looks better only, subjective-only acceptance
