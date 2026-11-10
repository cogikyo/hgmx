import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const gearGroup = el.querySelector("g.gear");
  animate(gearGroup, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmax });

  const head = el.querySelector('circle[cx="9"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });

  const body = el.querySelector('path[d*="15H6"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
