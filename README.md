## 💡 Features (What Can It Do?)

This awesome chat app will (eventually!) let you:

* Send and receive real-time messages.
* Create and join different chat rooms/channels.
* See who else is online.

## 🤝 Contributing (How You Can Help - Maybe Later!)

Right now, this project is being built by Debreczeni Alex János.
Thanks for your interest in contributing to this project! At this time, I'm focusing on developing it myself and not actively accepting external contributions. This is a personal project I'm using to learn and experiment, so I'm keeping the scope focused for now. I may consider accepting contributions like bug reports or feature requests in the future. If that changes, I'll be sure to update this file!

## 📄 License

MIT License

Copyright (c) 2025 Debreczeni Alex János

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## 🎉 Show Your Support!

## ✨ Tech Stack

This chat app is will powered by some really neat tools:

* **Frontend:** [Solid.js](https://www.solidjs.com/) - A super fast and reactive JavaScript library for building user interfaces. Think of it like a really efficient engine for making your website interactive!
* **Styling:** [Tailwind CSS](https://tailwindcss.com/) - A utility-first CSS framework that lets you style your app quickly by using pre-made classes. It's like having a big box of ready-to-use design Lego bricks!
* **Backend:** [Go](https://go.dev/) - A powerful and efficient programming language perfect for handling all the behind-the-scenes stuff, like managing messages and users. It's like the strong foundation of your app!
* **Database:** [Drizzle ORM](https://orm.drizzle.team/) - A modern and type-safe way to interact with your database. It helps your Go backend talk to and store information smoothly. Think of it as the librarian that organizes all your chat data!

## 🚀 Getting Started (For Developers!)

Want to see this in action or even help build it? Here's how you can get started on your own computer:

### Prerequisites

Make sure you have these installed on your system:

* [Go](https://go.dev/doc/install) (for the backend)
* [Bun](https://bun.sh/docs/installation) (for the frontend - it comes with its own package manager!)

### Running the Application

Follow these steps to get the chat app up and running:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/413ksz/BlueFox.git
    cd BlueFox
    ```

2.  **Set up the Backend:**
    ```bash
    cd backend
    go mod tidy  # Download any necessary Go packages
    # You might need to set up your database connection here. Check for a config file!
    go run main.go
    ```
    *(This will start your Go backend server. Keep this running in a separate terminal window!)*

3.  **Set up the Frontend:**
    ```bash
    cd frontend
    bun install       # Install all the frontend dependencies (like Solid.js and Tailwind) using Bun!
    bun run dev       # Start the frontend development server using Bun!
    ```
    *(This will usually open your chat app in your web browser, often at `http://localhost:3000` or something similar!)*

If you think this project is cool, feel free to give it a star ⭐ on GitHub! It helps others discover it too!
