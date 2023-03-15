import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import Swal from 'sweetalert2'
import Input from "./form/input";
import Select from "./form/Select";
import TextArea from "./form/TextArea";
import CheckBox from "./form/Checkbox";

const EditMovie = () => {
    const navigate = useNavigate();
    const { jwtToken } = useOutletContext();
    
    const [error, setError] = useState(null);
    //hold every error
    const [errors, setErrors] = useState([]);

    const hasError = (key) => {
        return errors.indexOf(key) !== -1;
    }

    const mpaaOptions = [
        {id: "G", value: "G"},
        {id: "PG", value: "PG"},
        {id: "PG13", value: "PG13"},
        {id: "R", value: "R"},
        {id: "NC17", value: "NC17"},
        {id: "18A", value: "18A"},
    ]

    const [movie, setMovie] = useState({
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        description: "",
        genres: [],
        genres_array: [Array(13).fill(false)],
    })

    //get id from the URL
    let {id} = useParams();
    if (id === undefined) {
        id = 0;
    }

    useEffect(()=> {
        if (jwtToken === "") {
            navigate("/login");
            return
        }

        // determine whether we are editing or adding movie
        if (id === 0) {
            //adding a movie
            setMovie({
                id: 0,
                title: "",
                release_date: "",
                runtime: "",
                mpaa_rating: "",
                description: "",
                genres: [],
                genres_array: [Array(13).fill(false)],
            })

            const headers = new Headers();
            headers.append("Content-Type", "application/json");

            const requestOptions = {
                method: "GET",
                headers: headers
            }

            fetch(`/genres`, requestOptions)
                .then(response => response.json())
                .then(data => {
                    const checks = [];
                    data.forEach(g => {
                        checks.push({id: g.id, checked: false, genre: g.genre})
                    })
                    
                    setMovie(prev => ({
                        ...prev, 
                        genres: checks,
                        genres_array: []
                    }))
                })
                .catch(err => {
                    console.log(err);
                })
        } else {
            // populate existing data to edit movie
            const headers = new Headers();
            headers.append("Content-Type", "application/json");
            headers.append("Authorization", "Bearer "+ jwtToken);

            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            fetch(`/admin/movies/${id}`, requestOptions)
                .then(response => {
                    if (response.status !== 200) {
                        setError("Invalid response code: "+response.status)
                    }
                    return response.json()
                })
                .then(data => {
                    //convert data
                    data.movie.release_date = new Date(data.movie.release_date).toISOString().split("T")[0]
                    const checks = [];

                    data.genres.forEach(g => {
                        if (data.movie.genres_array.indexOf(g.id) !== -1) {
                            checks.push({id: g.id, checked: true, genre: g.genre});
                        } else {
                            checks.push({id: g.id, checked: false, genre: g.genre});
                        }
                    })

                    setMovie({
                        ...data.movie,
                        genres: checks
                    })
                })
                .catch(err => {
                    console.log(err);
                })
        }

    }, [id, jwtToken, navigate])

    const handleSubmit = (event) =>{
        event.preventDefault();

        let errors = [];
        let required = [
            {field: movie.title, name: "title"},
            {field: movie.release_date, name: "release_date"},
            {field: movie.runtime, name: "runtime"},
            {field: movie.description, name: "description"},
            {field: movie.mpaa_rating, name: "mpaa_rating"},
        ]

        required.forEach(function (obj) {
            if (obj.field === "") {
                errors.push(obj.name);
            }
        })

        if (movie.genres_array.length === 0) {
            Swal.fire({
                title: "Error!",
                text: "You must choose at least one genre",
                icon: 'error',
                confirmButtonText: "OK",
            })
            errors.push("genres");
        }

        setErrors(errors);

        if (errors.length > 0){
            return false;
        }

        // verified => save changes
        const headers = new Headers();
        headers.append("Content-Type", "application/json");
        headers.append("Authorization", "Bearer "+jwtToken);

        // assume we are adding new movie
        let method = "PUT";

        // if we are editing existing movie, change to PATCH
        if (movie.id > 0) {
            method = "PATCH";
        }

        const requestBody = movie;
        // we need to convert string var from form

        requestBody.release_date = new Date(movie.release_date);
        requestBody.runtime = parseInt(movie.runtime, 10);

        let requestOptions = {
            body: JSON.stringify(requestBody),
            method: method,
            headers: headers,
            credentials: "include",
        }

        fetch(`/admin/movies/${movie.id}`, requestOptions)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    console.log(data.error); 
                } else {
                    navigate("/manage-catalogue")
                }
            })
            .catch(err => {
                console.log(err);
            })
    }

    const handleChange = (event) => {
        let value = event.target.value
        let name = event.target.name
        console.log(event.target);
        setMovie((prev)=> {
            return {
                ...prev,
                [name]: value
            }
        })
    }

    const handleCheck = (event, index) => {
        let tmpArr = movie.genres;
        tmpArr[index].checked = !tmpArr[index].checked;

        let tmpIDs = movie.genres_array;
        if (!event.target.checked) {
            tmpIDs.splice(tmpIDs.indexOf(event.target.value));
        } else {
            tmpIDs.push(parseInt(event.target.value, 10));
        }

        setMovie({
            ...movie,
            genres_array: tmpIDs,
        })
    }

    const confirmDelete = () => {
        Swal.fire({
            title: 'Delete movie?',
            text: "You cannot undo this action",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, delete it!'
          }).then((result) => {
            if (result.isConfirmed) {
                let headers = new Headers();
                headers.append("Authorization", "Bearer "+jwtToken)

                const requestOptions = {
                    method:"DELETE",
                    headers: headers,
                }
                fetch(`/admin/movies/${movie.id}`, requestOptions)
                    .then(response => response.json())
                    .then(data => {
                        if (data.error) {
                            console.log(data.error);
                        } else {
                            navigate("/manage-catalogue")
                        }
                    })
            }
          })
    }

    if (error !== null) {
        return <div>Error: {error.message}</div>
    } else {
    return(
        <div>
            <h2>Add/ Edit Movie</h2>
            <hr />
            {/* <pre>{JSON.stringify(movie, null, 3)}</pre> */}

            <form onSubmit={handleSubmit} action="">

                <input type="hidden" name="id" value={movie.id} id="id"></input>

                <Input
                    title={"Title"}
                    className={"form-control"}
                    type={"text"}
                    name={"title"}
                    value={movie.title}
                    onChange={handleChange}
                    errorDiv={hasError("title") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a title"}
                />

                <Input
                    title={"Release Date"}
                    className={"form-control"}
                    type={"date"}
                    name={"release_date"}
                    value={movie.release_date}
                    onChange={handleChange}
                    errorDiv={hasError("release_date") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a release date"}
                />

                <Input
                    title={"Runtime"}
                    className={"form-control"}
                    type={"text"}
                    name={"runtime"}
                    value={movie.runtime}
                    onChange={handleChange}
                    errorDiv={hasError("runtime") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a runtime"}
                />

                <Select 
                    title={"MPAA Rating"}
                    name={"mpaa_rating"}
                    options={mpaaOptions}
                    value={movie.mpaa_rating}
                    onChange={handleChange}
                    placeHolder={"Choose..."}
                    errorDiv={hasError("mpaa_rating") ? "text-danger" : "d-none"}
                    errorMsg={"Please choose"}
                />

                <TextArea
                    title={"Description"}
                    className={"form-control"}
                    type={"text"}
                    name={"description"}
                    value={movie.description}
                    rows={"3"}
                    onChange={handleChange}
                    errorDiv={hasError("description") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a description"}
                />
                <hr />

                <h3>Genres</h3>

                {movie.genres && movie.genres.length > 1 &&
                    <>
                        {Array.from(movie.genres).map((g, index) =>
                            <CheckBox
                                title={g.genre}
                                name={"genre"}
                                key={index}
                                id={"genre-"+index}
                                onChange={(event) => handleCheck(event, index)}
                                value={g.id}
                                checked={movie.genres[index].checked}
                            />
                        )}
                    </>
                }

                <hr />

                <button className="btn btn-primary">Save</button>

                {movie.id > 0 &&
                    <a href="#!" className="btn btn-danger ms-2" onClick={confirmDelete}>
                        Delete Movie
                    </a>
                }

            </form>
        </div>
    )}
}

export default EditMovie;