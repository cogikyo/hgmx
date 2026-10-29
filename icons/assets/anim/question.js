import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const circle = el.querySelector("circle");
  animate(circle, { pathLength: [0, 1], rotate: [270, 90] }, { duration: τmid });

  const questionMark = el.querySelectorAll("path");
  animate(questionMark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmax });
}
