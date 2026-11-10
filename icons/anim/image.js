import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const ridge = el.querySelector("path");
  animate(ridge, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
