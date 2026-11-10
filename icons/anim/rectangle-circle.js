import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rect = el.querySelector("path");
  animate(rect, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid, delay: τmin });
}
