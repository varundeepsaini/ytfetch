# YouTube Videos Dashboard

A modern React dashboard for viewing YouTube videos fetched by the ytfetch API.

## Features

- View videos in a responsive grid layout
- Infinite scroll pagination
- Date range filtering
- Modern Material-UI design
- Real-time updates

## Prerequisites

- Node.js 14.x or later
- npm 6.x or later
- Running ytfetch API server on port 8080

## Installation

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm start
```

The dashboard will be available at http://localhost:3000

## Development

- The dashboard uses Material-UI for components
- Axios for API calls
- React hooks for state management
- Responsive design for all screen sizes

## Building for Production

To create a production build:

```bash
npm run build
```

The build files will be in the `build` directory.

## Environment Variables

Create a `.env` file in the root directory with:

```
REACT_APP_API_URL=http://localhost:8080
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 