import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const medal = el.querySelector("circle");
  animate(medal, { scale: [0.5, 1.25, 1] }, { duration: τmax });
  const ribbon = el.querySelector("path");
  animate(ribbon, { rotate: [0, -15, 10, 0] }, { duration: τmid, delay: τmin });
}
