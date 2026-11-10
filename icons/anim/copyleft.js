import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const letter = el.querySelector("path");
  animate(letter, { rotate: [0, -45, 15, 0] }, { duration: τmid, delay: τmin });
}
