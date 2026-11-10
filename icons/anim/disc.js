import { animate, τmax } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('circle[r="10"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 630] }, { duration: τmax });
}
