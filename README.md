# Space Pizza

A terminal-based arcade game where you collect pizzas while avoiding enemy aliens!

## Overview

Space Pizza is a fast-paced terminal game built in Go, utilizing the `tcell` library for terminal graphics. Navigate through space, collect pizzas to increase your score, and avoid the aliens who are constantly hunting you down!

## Features

- Smooth, real-time terminal-based gameplay
- Intelligent enemy AI that actively hunts the player
- Score tracking system
- Randomly generated obstacles
- Game over screen with restart option
- Non-blocking input handling using Go channels and goroutines

## Requirements

- Go 1.16 or higher
- Terminal with support for basic ASCII characters

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/space-pizza.git
   cd space-pizza
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the game:
   ```bash
   go build -o space-pizza
   ```

4. Run the game:
   ```bash
   ./space-pizza
   ```

## How to Play

- Use the **arrow keys** to move your character (X) around the screen
- Collect food items (o) to increase your score
- Avoid enemies (E) who will try to catch you
- If an enemy catches you, the game is over

### Controls

- **↑ ↓ ← →**: Move your character
- **R**: Restart the game (when game over)
- **C or ESC**: Exit the game

## Game Elements

- **X**: Your character (the pizza delivery person)
- **o**: Pizza (worth 10 points each)
- **E**: Enemy aliens (they follow you with 80% probability)
- **.**: Obstacles (you cannot move through these)

## Game Mechanics

### Enemy AI

The enemies in Space Pizza use a simple but effective AI:
- 80% of the time, enemies will move toward the player
- 20% of the time, enemies will move in a random direction
- Enemies prioritize reducing the larger axis distance first
- If an enemy cannot move in its chosen direction, it will pick a new random direction

### Game Loop

The game uses non-blocking input handling through Go channels to ensure smooth gameplay:
- A separate goroutine handles keyboard input
- The main game loop processes movement, collisions, and rendering
- Enemy movement and game updates occur on a timer regardless of player input

## Development

The game is built using Go and the `tcell` library for terminal graphics. The codebase is structured as follows:

- `game.go`: Main game code including game loop, input handling, and rendering
- Types:
  - `Game`: Main game state
  - `Player`: Player position and score
  - `Enemy`: Enemy position and movement data
  - `Food`: Food position data

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by classic arcade games
- Built with [tcell](https://github.com/gdamore/tcell) library 