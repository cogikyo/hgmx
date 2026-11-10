import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ear = el.querySelector("g");
  animate(ear, { rotate: [0, -8, 3, 0] }, { duration: τmid });

  const slash = el.querySelector("line");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
