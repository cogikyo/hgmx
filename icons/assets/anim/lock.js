import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shackle = el.querySelector("path");
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid });

  const body = el.querySelector("rect");
  animate(body, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
