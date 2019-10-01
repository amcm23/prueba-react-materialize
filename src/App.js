import React, { useEffect, useState } from 'react'
import M from 'materialize-css';
import { Navbar, NavItem, TextInput, Table } from 'react-materialize'
import Button from 'react-materialize/lib/Button';
import axios from 'axios'
import Swal from 'sweetalert2'
// ref can only be used on class components
export default function App() {
  // get a reference to the element after the component has mounted
  useEffect(
    () => {
      M.AutoInit();
    }
  )

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = () => {
    axios.get('http://localhost:8000/posts')
      .then(response => {
        console.log(response.data);
        setData(response.data)
      })
      .catch(error => {
        console.log(error);
      });
  }

  const [name, setName] = useState("")
  const [surname, setSurname] = useState("")
  const [data, setData] = useState()

  function handleNameChange(e) {
    setName(e.target.value)
  }
  function handleSurnameChange(e) {
    setSurname(e.target.value)
  }

  function handleSubmit() {
    const values = JSON.stringify({
      "nombre": name,
      "apellido": surname
    })

    console.log(values)

    axios({
      method: 'post',
      url: 'http://localhost:8000/posts',
      data: values,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    }).then(function (response) {
      console.log(response);
      fetchUsers()
      Swal.fire(({
        position: 'top',
        type: 'success',
        title: 'Usuario añadido con éxito.',
        showConfirmButton: false,
        timer: 1500
      }))
      setName("")
      setSurname("")


    }).catch(function (error) {
      console.log(error);
    });
  }

  const deleteUser = (id) => {

    Swal.fire({
      title: '¿Está seguro que desea eliminar el usuario?',
      text: "Esta acción será irreversible.",
      type: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: 'Si, borralo.'
    }).then((result) => {
      if (result.value) {
        axios({
          method: 'delete',
          url: `http://localhost:8000/posts/${id}`,
          data: null,
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            'Accept': 'application/json'
          }
        }).then(response => {
          console.log(response.data);
          Swal.fire(
            '¡Eliminado!',
            'El usuario ha sido eliminado con éxito.',
            'success'
          )
          fetchUsers()
        })
          .catch(error => {
            console.log(error);
          });
      }
    })
  }

  return (
    <div>
      <Navbar brand={<a>PruebaMaterialize</a>} alignLinks="right">
        <NavItem href="">
          Getting started
          </NavItem>
        <NavItem href="components.html">
          Components
          </NavItem>
      </Navbar>
      <div style={{ padding: "2rem" }}>
        <h1>Ejemplo de formulario</h1>
        <TextInput label="Nombre" value={name} onChange={handleNameChange} />

        <TextInput label="Apellidos" value={surname} onChange={handleSurnameChange} />

        <Button onClick={handleSubmit}>Registrar</Button>
        <hr style={{ marginTop: "3rem" }} />
        <Table style={{ marginTop: "1rem" }}>
          <thead>
            <tr>
              <th data-field="id">
                ID
              </th>
              <th data-field="name">
                Nombre
              </th>
              <th data-field="surname">
                Apellidos
              </th>
              <th>
                Acciones
              </th>
            </tr>
          </thead>
          <tbody>
            {data && data.map((user, i) => (
              <tr key={i}>
                <td>
                  {user.id}
                </td>
                <td>
                  {user.nombre}
                </td>
                <td>
                  {user.apellido}
                </td>
                <td>
                  <Button onClick={() => deleteUser(user.id)}>Eliminar</Button>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </div>
    </div>
  )
}
