import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const barn = el.querySelectorAll('path[d^="M22 12H2"], path[d^="M21 12v8"]');
  animate(barn, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector('path[d^="m7 22"]');
  animate(door, { scaleY: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}
