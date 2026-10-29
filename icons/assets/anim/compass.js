import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const needle = el.querySelector("path");
  animate(needle, { rotate: [0, 390, 345, 360] }, { duration: τmid, delay: τmin });
}
