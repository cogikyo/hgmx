import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const [upper, lower] = el.querySelectorAll("g");
  animate(upper, { rotate: [0, 15, -5, 0] }, { duration: τmid, delay: τmin });
  animate(lower, { rotate: [0, -15, 5, 0] }, { duration: τmid, delay: τmin });
}
