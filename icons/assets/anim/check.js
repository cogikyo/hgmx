import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tick = el.querySelector("path");
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
