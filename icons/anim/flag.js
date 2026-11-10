import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const flag = el.querySelector("path");
  animate(flag, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(flag, { rotate: [0, -10, 6, -3, 0] }, { duration: τmid, delay: τmin });
}
