import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M4 10.95"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const awning = el.querySelector('path[d^="M17.774"]');
  animate(awning, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const door = el.querySelector('path[d^="M15 21"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
