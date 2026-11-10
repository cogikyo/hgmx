import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const slash = el.querySelector("line");
  animate(slash, { scale: [1, 1.25, 1] }, { duration: τmid, delay: τmin });
}
