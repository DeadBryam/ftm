# TUI Redesign: RPG Minimal

## Overview
Redesign the terminal UI to be beautiful, intuitive, and friendly for non-technical D&D Dungeon Masters while maintaining a clean, minimal aesthetic.

## Goals
- Make the TUI as beautiful and usable as the web dashboard
- Help non-technical users understand the tool intuitively
- Use thematic RPG styling without over-decoration
- Provide clear visual feedback for tunnel states
- Implement a 2-column layout for large terminals

## Design: "RPG Minimal"

### Color Palette
| Role | Color | Hex | Usage |
|------|-------|-----|-------|
| Background | Dark Parchment | `#1a1814` | Terminal background |
| Primary Accent | Antique Gold | `#c9a227` | Selection, titles, highlights |
| Secondary Accent | Aged Bronze | `#8b7355` | Subtle borders, dividers |
| Text Primary | Bone White | `#e8e6e1` | Main text |
| Text Secondary | Warm Gray | `#9a9590` | Hints, secondary info |

### Status Background Colors
| Status | Background | Usage |
|--------|-----------|-------|
| Online | `#1e3a2f` (Forest Green) | Running tunnels |
| Offline | `#2a2824` (Parchment Gray) | Stopped tunnels |
| Starting | `#3a3020` (Dark Amber) | Starting tunnels |
| Error | `#3a2020` (Wine Red) | Error state |

### Layout
Two-column responsive layout for terminals ≥100 columns:
- **Left Column (40%):** Tunnel list with visual state indicators
- **Right Column (60%):** Detail view for selected tunnel

For small terminals (<100 columns): Single column with compact list view.

### Visual Elements

#### Header
```
🎲  Foundry Tunnel Manager                    v0.2.0
```
- Gold title, right-aligned version
- No box decorations, clean text only

#### Tunnel List Items
Each tunnel displayed as:
```
🟢 Campaign Name              Online
   Provider • :30000
```
- Status emoji (🟢🟡🔴) at left
- Name in bold when selected
- Background color indicates status
- Selection: gold left border indicator

#### Detail Panel (Right Column)
Shows for selected tunnel:
- Full name and provider
- Status with colored indicator
- Public URL in highlighted box
- Action buttons: [Stop] [Logs] [Delete]

#### Empty State
Centered welcome screen for first-time users:
- Large decorative icon (🌐 + 🎲)
- Friendly welcome message
- Prominent CTA button
- Keyboard shortcut hint
- Web dashboard tip

#### Help Bar
Context-aware help at bottom:
```
↑/↓ navigate • Enter toggle • 'a' add • 'l' logs • 'q' quit
```

### Animations
- Selection change: instant (no animation needed in terminal)
- Status transitions: color fade via periodic updates

### Keyboard Navigation
| Key | Action |
|-----|--------|
| ↑/↓ | Navigate tunnels |
| Enter | Toggle selected tunnel on/off |
| a | Add new tunnel |
| l | View logs for selected |
| d | Delete selected |
| e | Edit selected |
| c | Copy URL to clipboard |
| w | Open web dashboard |
| q/ctrl+c | Quit |

## Technical Implementation

### Components Needed
1. **LayoutManager** - Handles responsive 1/2 column layout based on terminal width
2. **TunnelList** - Render list with status backgrounds and selection
3. **DetailPanel** - Right column detail view with actions
4. **EmptyState** - Welcome screen for no tunnels
5. **HelpBar** - Contextual keyboard help
6. **Styled URL Box** - Highlighted URL with copy action

### Files to Modify
- `internal/app/view.go` - Main view functions
- `internal/app/model.go` - Add new state fields if needed
- `internal/app/styles.go` - Define all color styles

### Dependencies
- Existing: `charmbracelet/lipgloss`, `charmbracelet/bubbletea`
- No new dependencies required

## Success Criteria
- [ ] TUI is visually appealing with RPG minimal theme
- [ ] 2-column layout works in terminals ≥100 columns
- [ ] Single-column fallback for small terminals
- [ ] Empty state welcomes first-time users
- [ ] Tunnel states are visually distinct (background colors)
- [ ] Detail panel shows all relevant info + actions
- [ ] Help bar is contextual and helpful
- [ ] All existing keyboard shortcuts still work

## Notes
- Keep existing update mechanism - no changes to model logic
- Focus on View() functions and styling only
- Maintain compatibility with existing commands
- No comments in code per project convention
