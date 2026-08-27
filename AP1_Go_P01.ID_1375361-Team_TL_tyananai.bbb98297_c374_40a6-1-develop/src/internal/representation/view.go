package representation

import (
	"fmt"

	"rogue/internal/datalayer"
	"rogue/internal/domain"

	"github.com/nsf/termbox-go"
)

const savePath = datalayer.SavePath

const (
	SCREEN_MENU = iota
	SCREEN_GAME
	SCREEN_FOODS
	SCREEN_ELIXIRS
	SCREEN_ROLLS
	SCREEN_WEAPONS
	SCREEN_RESULT
	SCREEN_EXIT
	SCREEN_NO_SAVE
	SCREEN_NAME_INPUT
	SCREEN_LEADERS
	SCREEN_ERROR
	START_NEW_GAME_MENU_ITEM = "Начать новую игру"
	CONTINUE_GAME_MENU_ITEM  = "Продолжить игру"
	RECORDS_MENU_ITEM        = "Достижения"
	EXIT_MENU_ITEM           = "Выйти"
)

type View struct {
	Game       *domain.Game
	State      int
	menu       Menu
	started    bool
	playerName string
	nameInput  string
	errorMsg   string
}

type Menu struct {
	Values []string
	CurIdx int
}

func NewView() *View {
	return &View{
		Game:  nil,
		State: SCREEN_MENU,
		menu: Menu{
			Values: []string{
				START_NEW_GAME_MENU_ITEM,
				CONTINUE_GAME_MENU_ITEM,
				RECORDS_MENU_ITEM,
				EXIT_MENU_ITEM,
			},
			CurIdx: 0,
		},
		started: false,
	}
}

func (v *View) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	switch v.State {
	case SCREEN_MENU:
		v.renderMenu()
	case SCREEN_GAME:
		v.renderGame()
	case SCREEN_FOODS:
		v.renderFoods()
	case SCREEN_ELIXIRS:
		v.renderElixirs()
	case SCREEN_WEAPONS:
		v.renderWeapons()
	case SCREEN_ROLLS:
		v.renderRolls()
	case SCREEN_RESULT:
		v.renderResults()
	case SCREEN_NO_SAVE:
		v.renderNoSave()
	case SCREEN_NAME_INPUT:
		v.renderNameInput()
	case SCREEN_LEADERS:
		v.renderLeaders()
	case SCREEN_ERROR:
		v.renderError()
	}
	termbox.Flush()
}

func (v *View) renderMenu() {
	for i, str := range v.menu.Values {
		v.drawString("   "+str, 3, i)
	}
	v.drawString(" > ", 0, v.menu.CurIdx)
}

func (v *View) renderGame() {
	lm := v.Game.CurrentLevelMap
	player := v.Game.Player

	for y := range domain.HEIGHT {
		for x := range domain.WIDTH {
			cell := lm.Field[y][x]
			background := termbox.ColorDefault
			if cell.Visibility {
				switch cell.EnvType {
				case domain.EnvFood:
					background = termbox.ColorGreen
				case domain.EnvElixir:
					background = termbox.ColorMagenta
				case domain.EnvRoll:
					background = termbox.ColorLightYellow
				case domain.EnvWeapon:
					background = termbox.ColorRed
				}
				termbox.SetCell(x, y, cell.EnvType, background, termbox.ColorDefault)
			}
		}
	}

	for _, e := range lm.Enemies {
		if !e.Alive || !e.IsVisible || !lm.Field[e.Pos.Y][e.Pos.X].Visibility {
			continue
		}
		termbox.SetCell(e.Pos.X, e.Pos.Y, e.Symbol, termbox.ColorDefault, termbox.ColorDefault)
	}

	kP := lm.Key.PosOnMap
	if lm.Field[kP.Y][kP.X].Visibility {
		termbox.SetCell(kP.X, kP.Y, '*', termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(player.Pos.X, player.Pos.Y, player.Symbol, termbox.ColorDefault, termbox.ColorDefault)

	weaponStr := "нет"
	if player.CurrentWeapon != nil {
		weaponStr = fmt.Sprintf("сила +%d", player.CurrentWeapon.Power)
	}

	elixirStr := ""
	if player.ElixirMovesLeft > 0 {
		elixirStr = fmt.Sprintf("  Эликсир: %d ходов", player.ElixirMovesLeft)
	}

	personStat := fmt.Sprintf("Уровень: %d/21  Здоровье:%d(%d)  Оружие: %s%s  Золото: %d",
		lm.Level, player.Health, player.MaxHealth, weaponStr, elixirStr, player.Backpack.Gold)
	v.drawString(personStat, 0, 28)

	backpackStat := fmt.Sprintf("Рюкзак - Еда:%d  Эликсир:%d  Свитки:%d  Оружие:%d", len(player.Backpack.Food), len(player.Backpack.Elixir), len(player.Backpack.Roll), len(player.Backpack.Weapon))
	v.drawString(backpackStat, 0, 29)

	hints := "W/A/S/D - движение | H - оружие | J - еда | K - эликсир | E - свитки | Q - выход"
	v.drawString(hints, 0, 30)
	task := "Найди ключ * чтобы перейти на следующий уровень"
	v.drawString(task, 0, 31)
	if v.Game.LastLog != "" {
		v.drawString(v.Game.LastLog, 0, 32)
	}
}

func (v *View) drawString(s string, x int, y int) {
	xi := x
	for _, ch := range s {
		termbox.SetCell(xi, y, ch, termbox.ColorDefault, termbox.ColorDefault)
		xi++
	}
}

func (v *View) renderFoods() {
	v.drawString("--- Еда ---", 0, 0)
	food := v.Game.Player.Backpack.Food
	if len(food) == 0 {
		v.drawString("Еда отсутствует в рюкзаке (q - вернуться на карту)", 0, 2)
	} else {
		hint := fmt.Sprintf("Выбрать (1-%d, q - вернуться на карту): ", len(food))
		v.drawString(hint, 0, 2)
		for i, f := range food {
			foodItem := fmt.Sprintf("%d. Еда (+%d HP)", i+1, f.HP)
			v.drawString(foodItem, 0, i+3)
		}
	}
}

func (v *View) renderElixirs() {
	elixirs := v.Game.Player.Backpack.Elixir
	v.drawString("--- Эликсиры ---", 0, 0)

	if len(elixirs) == 0 {
		v.drawString("Эликсиры отсутствуют в рюкзаке (q - вернуться)", 0, 2)
		return
	} else {
		hint := fmt.Sprintf("Выбрать (1-%d, q - вернуться на карту): ", len(elixirs))
		v.drawString(hint, 0, 2)
		for i, e := range elixirs {
			effect := elixirEffectStr(e)
			elexirStr := fmt.Sprintf("%d. Эликсир (%s, %d ходов)", i+1, effect, e.NumMoves)
			v.drawString(elexirStr, 0, 3+i)
		}
	}
}

func elixirEffectStr(e domain.Elixir) string {
	if e.MaxHP > 0 {
		return fmt.Sprintf("+%d макс. HP", e.MaxHP)
	}
	if e.Agility > 0 {
		return fmt.Sprintf("+%d ловкость", e.Agility)
	}
	if e.Power > 0 {
		return fmt.Sprintf("+%d сила", e.Power)
	}
	return "неизвестный эффект"
}

func (v *View) renderWeapons() {
	weapons := v.Game.Player.Backpack.Weapon
	curWeapon := v.Game.Player.CurrentWeapon
	v.drawString("--- Оружие ---", 0, 0)
	if curWeapon == nil {
		v.drawString("Оружия в руках нет", 0, 2)
	} else {
		curWeaponStr := fmt.Sprintf("Оружие в руках: (сила +%d, убрать - 0)", curWeapon.Power)
		v.drawString(curWeaponStr, 0, 2)
	}
	if len(weapons) == 0 {
		v.drawString("Оружия в рюкзаке нет (q - вернуться на карту)", 0, 3)
	} else {
		hint := fmt.Sprintf("Выбрать (1-%d, q - вернуться на карту): ", len(weapons))
		v.drawString(hint, 0, 3)
		for i, w := range weapons {
			weaponStr := fmt.Sprintf("%d. Оружие (сила +%d)", i+1, w.Power)
			v.drawString(weaponStr, 0, 5+i)
		}
	}
}

func (v *View) renderRolls() {
	v.drawString("--- Свитки ---", 0, 0)
	rolls := v.Game.Player.Backpack.Roll
	if len(rolls) == 0 {
		v.drawString("Свитки отсутствуют в рюкзаке (q - вернуться на карту)", 0, 2)
	} else {
		hint := fmt.Sprintf("Что выбрать (1-%d, q - вернуться на карту): ", len(rolls))
		v.drawString(hint, 0, 2)
		for i, r := range rolls {
			effect := rollEffectStr(r)
			rollStr := fmt.Sprintf("%d. Свиток (%s, постоянно", i+1, effect)
			v.drawString(rollStr, 0, 3+i)
		}
	}
}

func rollEffectStr(r domain.Roll) string {
	if r.MaxHP > 0 {
		return fmt.Sprintf("+%d макс. HP", r.MaxHP)
	}
	if r.Agility > 0 {
		return fmt.Sprintf("+%d ловкость", r.Agility)
	}
	if r.Power > 0 {
		return fmt.Sprintf("+%d сила", r.Power)
	}
	return "неизвестный эффект"
}

func (v *View) renderResults() {
	v.drawString("GAME OVER", 0, 0)
	v.drawString("press any key to close", 0, 1)
}

func (v *View) renderNoSave() {
	v.drawString("Нет сохраненных данных или они некорректны", 0, 0)
	v.drawString("q - вернуться в главное меню", 0, 1)
}

func (v *View) renderError() {
	v.drawString("Ошибка: "+v.errorMsg, 0, 0)
	v.drawString("q - вернуться в главное меню", 0, 2)
}

func (v *View) renderLeaders() {
	entries := datalayer.LoadLeaders(datalayer.LeadersPath)
	v.drawString("Таблица лидеров:", 0, 0)
	if len(entries) == 0 {
		v.drawString("Пока нет записей", 0, 2)
	} else {
		for i, e := range entries {
			line := fmt.Sprintf("%d. %s - %d уровень - %d золота", i+1, e.Name, e.Level, e.Gold)
			v.drawString(line, 0, i+2)
		}
	}
	v.drawString("q - вернуться в главное меню", 0, 13)
}

func (v *View) HandleInput() termbox.Event {
	e := termbox.PollEvent()
	for e.Type != termbox.EventKey {
		e = termbox.PollEvent()
	}
	return e
}

func (v *View) Update(e termbox.Event) {
	switch v.State {
	case SCREEN_MENU:
		v.State = v.handleMenu(e)
	case SCREEN_GAME:
		v.State = v.handleGame(e)
	case SCREEN_FOODS:
		v.State = v.handleFoods(e)
	case SCREEN_ELIXIRS:
		v.State = v.handleElixirs(e)
	case SCREEN_WEAPONS:
		v.State = v.handleWeapons(e)
	case SCREEN_ROLLS:
		v.State = v.handleRolls(e)
	case SCREEN_RESULT:
		v.State = v.handleResult(e)
	case SCREEN_NO_SAVE:
		if e.Ch == 'q' {
			v.State = SCREEN_MENU
		}
	case SCREEN_NAME_INPUT:
		v.State = v.handleNameInput(e)
	case SCREEN_LEADERS:
		if e.Ch == 'q' {
			v.State = SCREEN_MENU
		}
	case SCREEN_ERROR:
		if e.Ch == 'q' {
			v.State = SCREEN_MENU
		}
	}
}

func (v *View) handleMenu(e termbox.Event) int {
	switch {
	case e.Key == termbox.KeyEnter:
		switch v.menu.CurIdx {
		case 0:
			v.nameInput = ""
			return SCREEN_NAME_INPUT
		case 1:
			g, err := datalayer.LoadGame(savePath)
			if err != nil {
				return SCREEN_NO_SAVE
			}
			v.Game = g
			return SCREEN_GAME
		case 2:
			return SCREEN_LEADERS
		case 3:
			return SCREEN_EXIT
		}
	case e.Ch == 's':
		v.menu.CurIdx = min(len(v.menu.Values)-1, v.menu.CurIdx+1)
	case e.Ch == 'w':
		v.menu.CurIdx = max(0, v.menu.CurIdx-1)
	}
	return SCREEN_MENU
}

func (v *View) handleGame(e termbox.Event) int {
	switch ch := e.Ch; ch {
	case 'q':
		return SCREEN_MENU
	case 'h':
		return SCREEN_WEAPONS
	case 'j':
		return SCREEN_FOODS
	case 'k':
		return SCREEN_ELIXIRS
	case 'e':
		return SCREEN_ROLLS
	default:
		prevLevel := v.Game.CurrentLevelMap.Level
		v.Game.Update(string(ch))
		if v.Game.PlayerDead {
			if err := datalayer.SaveLeader(v.playerName, v.Game.CurrentLevelMap.Level, v.Game.Player.Backpack.Gold); err != nil {
				return v.showError(err)
			}
			return SCREEN_RESULT
		}
		if v.Game.CurrentLevelMap.Level != prevLevel {
			if err := datalayer.SaveGame(v.Game, savePath); err != nil {
				return v.showError(err)
			}
			if err := datalayer.SaveLeader(v.playerName, v.Game.CurrentLevelMap.Level, v.Game.Player.Backpack.Gold); err != nil {
				return v.showError(err)
			}
		}
	}
	return SCREEN_GAME
}

func (v *View) handleFoods(e termbox.Event) int {
	if e.Ch == 'q' {
		return SCREEN_GAME
	}
	idx := int(e.Ch - '0')
	foods := v.Game.Player.Backpack.Food
	if idx >= 1 && idx <= len(foods) {
		domain.UseFood(v.Game.Player, idx)
	}
	return SCREEN_FOODS
}

func (v *View) handleElixirs(e termbox.Event) int {
	if e.Ch == 'q' {
		return SCREEN_GAME
	}
	idx := int(e.Ch - '0')
	elixirs := v.Game.Player.Backpack.Elixir
	if idx >= 1 && idx <= len(elixirs) {
		domain.UseElixir(v.Game.Player, idx)
	}
	return SCREEN_ELIXIRS
}

func (v *View) handleWeapons(e termbox.Event) int {
	if e.Ch == 'q' {
		return SCREEN_GAME
	}
	idx := int(e.Ch - '0')
	weapons := v.Game.Player.Backpack.Weapon
	if idx >= 0 && idx <= len(weapons) {
		domain.EquipWeapon(v.Game.Player, v.Game.CurrentLevelMap, idx)
	}
	return SCREEN_WEAPONS
}

func (v *View) handleRolls(e termbox.Event) int {
	if e.Ch == 'q' {
		return SCREEN_GAME
	}
	idx := int(e.Ch - '0')
	rolls := v.Game.Player.Backpack.Roll
	if idx >= 0 && idx <= len(rolls) {
		domain.UseRoll(v.Game.Player, idx)
	}
	return SCREEN_ROLLS
}

func (v *View) IsGameOver() bool {
	return v.State == SCREEN_EXIT
}

func (v *View) handleResult(_ termbox.Event) int {
	return SCREEN_EXIT
}

func (v *View) renderNameInput() {
	v.drawString("Введите ваше имя: "+v.nameInput, 0, 0)
	v.drawString("Enter - начать игру", 0, 2)
}

func (v *View) handleNameInput(e termbox.Event) int {
	switch {
	case e.Key == termbox.KeyEnter:
		if v.nameInput == "" {
			return SCREEN_NAME_INPUT
		}
		v.playerName = v.nameInput
		v.Game = domain.NewGame()
		if err := datalayer.SaveGame(v.Game, savePath); err != nil {
			return v.showError(err)
		}
		return SCREEN_GAME
	case e.Key == termbox.KeyBackspace || e.Key == termbox.KeyBackspace2:
		runes := []rune(v.nameInput)
		if len(runes) > 0 {
			v.nameInput = string(runes[:len(runes)-1])
		}
	case e.Ch != 0:
		v.nameInput += string(e.Ch)
	}
	return SCREEN_NAME_INPUT
}

func (v *View) showError(err error) int {
	v.errorMsg = err.Error()
	return SCREEN_ERROR
}
