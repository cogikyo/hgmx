import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const fans = [el.querySelector('circle[cx="8"]'), el.querySelector('circle[cx="16"]')];
  animate(fans, { rotate: [0, 360], pathLength: [0.25, 1] }, { duration: τmid, delay: stagger(0.12) });
}
