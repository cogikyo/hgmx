import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmax });

  const head = el.querySelector('circle[cx="10"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });

  const body = el.querySelector('path[d^="M2 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
