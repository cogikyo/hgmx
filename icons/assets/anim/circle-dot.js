import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('circle[r="10"]');
  animate(ring, { scale: [1, 0.85, 1] }, { duration: τmid });
  const dot = el.querySelector('circle[r="1"]');
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
