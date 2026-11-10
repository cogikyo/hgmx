import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const ground = el.querySelector('path[d="M2 21h20"]');
  animate(ground, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector("g");
  animate(door, { scaleX: [1, 0.8, 1] }, { duration: τmax });
}
