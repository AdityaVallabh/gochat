import { useEffect, useState } from 'react';

export default function Home() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const u = localStorage.getItem('user');
    setUser(JSON.parse(u));
  }, []);

  return (
    <>
      <h1>Welcome to GoChat!</h1>
      <div>{user ? <UserView user={user} /> : <NewUser />}</div>
    </>
  );
}

export function NewUser() {
  async function handleCreateUser(e) {
    e.preventDefault();
    const resp = await fetch('http://localhost:8000/api/user', {
      method: 'post',
      body: JSON.stringify({
        name: e.target.name.value,
      }),
    });
    const user = await resp.json();
    localStorage.setItem('user', JSON.stringify(user));
  }
  return (
    <div>
      <form onSubmit={handleCreateUser}>
        <label htmlFor="name">Who are you?</label>
        <input type="text" id="name" />
        <br />
      </form>
    </div>
  );
}

export function UserView({ user }) {
  const [room, setRoom] = useState(null);
  const [ws, setWs] = useState(null);

  async function handleRoomCreate() {
    const resp = await fetch('http://localhost:8000/api/room', {
      method: 'post',
    });
    const r = await resp.json();
    console.log(r);
    setRoom(r);
    join(r.ID, user.ID);
  }

  async function handleRoomJoin(e) {
    e.preventDefault();
    join(e.target.roomID.value, user.ID)
  }

  async function join(roomID, userID) {
    const resp = await fetch('http://localhost:8000/api/room/join', {
      method: 'post',
      body: JSON.stringify({
        RoomID: roomID,
        UserID: userID,
      }),
    });
    setRoom(await resp.json());
  }

  useEffect(() => {
    if (ws) {
      return;
    }
    const socket = new WebSocket('ws://localhost:8000/ws?id=' + user.ID);
    console.log('Attempting Connection...');

    socket.onopen = () => {
      console.log('Successfully Connected');
    };

    socket.onclose = (event) => {
      console.log('Socket Closed Connection: ', event);
      socket.send('Client Closed!');
    };

    socket.onerror = (error) => {
      console.log('Socket Error: ', error);
    };

    setWs(socket);
  }, []);

  return (
    <div>
      <p>Hello {user.Name}</p>
      {room ? (
        <RoomView user={user} room={room} ws={ws} />
      ) : (
        <div>
          <button onClick={handleRoomCreate}>Create Room</button>
          <form onSubmit={handleRoomJoin}>
            <input id="roomID" />
            <button>Join Room</button>
          </form>
        </div>
      )}
    </div>
  );
}

export function RoomView({ user, room, ws }) {
  const [messages, setMessages] = useState([]);
  const [count, setCount] = useState(0);
  useEffect(() => {
    ws.onmessage = (msg) => {
      console.log(msg.data)
      const mm = messages;
      const j = JSON.parse(msg.data);
      mm.push(j.Data);
      setMessages(mm);
      setCount(c => c+1);
      console.log(count);
    };
  }, []);
  function send(e) {
    e.preventDefault();
    ws.send(JSON.stringify({
      User: {
        ID: user.ID,
      },
      Room: {
        ID: room.ID,
      },
      Data: e.target.msg.value,
    }));
  }
  return (
    <div>
      {room.ID}
      {ws ? (
      <div>
        <form onSubmit={send}>
          <input id="msg" />
        </form>
        <div>
          {count}
          {messages.map(m => <p id={m}>{m}</p>)}
        </div>
      </div>) : null}
    </div>
  );
}
