import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, createRoutesFromElements, Route, RouterProvider } from 'react-router-dom';
import App from './App';
import EditMovie from './components/EditMovie';
import ErrorPage from './components/ErrorPage';
import Genres from './components/Genres';
import GraphQL from './components/GraphQL';
import Home from './components/Home';
import Login from './components/Login';
import ManageCatalogue from './components/ManageCatalogue';
import Movies from './components/Movies';
import Movie from './components/Movie';

// const router = createBrowserRouter([
//   {
//     path: "/",
//     element: <App />,
//     errorElement: <ErrorPage />,
//     children: [
//       {index: true, element: <Home /> },
//       {
//         path: "/movies",
//         element: <Movies />,
//       },
//       {
//         path: "/movies/:id",
//         element: <Movie />,
//       },
//       {
//         path: "/genres",
//         element: <Genres />,
//       },
//       {
//         path: "/admin/movie/0",
//         element: <EditMovie />,
//       },
//       {
//         path: "/manage-catalogue",
//         element: <ManageCatalogue />,
//       },
//       {
//         path: "/graphql",
//         element: <GraphQL />,
//       },
//       {
//         path: "/login",
//         element: <Login />,
//       },
//     ]
//   }
// ])

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path='/' element={<App />} errorElement={<ErrorPage />}>
      <Route index element={<Home />}></Route>
      <Route path='movies' element={<Movies />}></Route>
      <Route path='movies/:id' element={<Movie />}></Route>
      <Route path='genres' element={<Genres />}></Route>
      <Route path='admin/movie/0' element={<EditMovie />}></Route>
      <Route path='manage-catalogue' element={<ManageCatalogue />}></Route>
      <Route path='/graphql' element={<GraphQL />}></Route>
      <Route path='/login' element={<Login />}></Route>
    </Route>
  )
)

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
