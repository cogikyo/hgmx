import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const racks = el.querySelectorAll('path[d^="M4.5"]');
  animate(racks, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
