import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const key = el.querySelector("path");
  animate(key, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
  animate(key, { scale: [1, 0.85, 1] }, { duration: τmin, delay: τmid });
}
