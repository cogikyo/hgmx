import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const shackle = el.querySelector("path");
  animate(shackle, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const body = el.querySelector("rect");
  animate(body, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
