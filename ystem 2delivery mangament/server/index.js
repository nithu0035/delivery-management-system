/**
 * Minimal Express server scaffold (added by assistant).
 * Replace or expand with your existing server code as needed.
 */

require('dotenv').config();
const express = require('express');
const cors = require('cors');
const app = express();
app.use(cors());
app.use(express.json());

const PORT = process.env.PORT || 4000;

// Simple health route
app.get('/api/health', (req, res) => res.json({status: 'ok'}));

// Auth stub
app.post('/api/auth/login', (req, res) => {
  // NOTE: replace with real auth (DB lookup)
  const {email, password} = req.body;
  if(email === 'admin@example.com' && password === 'password') {
    const jwt = require('jsonwebtoken');
    const token = jwt.sign({email, role:'admin'}, process.env.JWT_SECRET || 'devsecret', {expiresIn:'7d'});
    return res.json({token});
  }
  return res.status(401).json({error:'invalid credentials'});
});

app.listen(PORT, () => console.log(`Server listening on ${PORT}`));
