import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M10.584 21"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector('path[d^="M14 12H10"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });

  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
