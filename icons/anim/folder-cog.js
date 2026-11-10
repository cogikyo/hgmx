import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M10.3"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
