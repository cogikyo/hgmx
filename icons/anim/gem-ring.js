import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });
  const gem = el.querySelector('path[d^="M13.2"]');
  animate(gem, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
