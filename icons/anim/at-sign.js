import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelector("circle");
  animate(bowl, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid });

  const tail = el.querySelector("path");
  animate(tail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
