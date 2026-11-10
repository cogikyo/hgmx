import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shield = el.querySelector('path[d^="M20 13"]');
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const gear = el.querySelector("g.gear");
  animate(gear, { rotate: [0, 360], scale: [1, 1.25, 1] }, { duration: τmid, delay: τmin });
}
