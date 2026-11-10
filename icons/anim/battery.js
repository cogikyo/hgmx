import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("rect");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const tip = el.querySelector("path");
  animate(tip, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
