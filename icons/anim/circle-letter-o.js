import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('circle[r="10"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const inner = el.querySelector('circle[r="4"]');
  animate(inner, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid, delay: τmin });
}
