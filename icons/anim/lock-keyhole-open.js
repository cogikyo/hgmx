import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const shackle = el.querySelector("g");
  animate(shackle, { rotate: [0, -20, 5, 0] }, { duration: τmax });

  const body = el.querySelectorAll("rect, circle");
  animate(body, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
