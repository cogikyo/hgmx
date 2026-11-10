import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector("rect");
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector('path[d^="M10 22"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
