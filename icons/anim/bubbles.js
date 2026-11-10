import { animate, stagger } from "../motion.js";

export default function (el) {
  const bubbles = [el.querySelector("g"), el.querySelector('circle[cx="18.5"]'), el.querySelector('circle[cy="4.5"]')];
  animate(bubbles, { scale: [1, 0.85, 1.1, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
