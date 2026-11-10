import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const mane = el.querySelector("path");
  animate(mane, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const head = el.querySelector("circle");
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
