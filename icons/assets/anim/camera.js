import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector("circle");
  animate(lens, { scale: [1, 0, 1.25, 1] }, { duration: τmax });
}
